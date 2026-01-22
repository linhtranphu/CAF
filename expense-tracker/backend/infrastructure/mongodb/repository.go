package mongodb

import (
	"context"
	"log"
	"os"
	"time"
	"expense-tracker/domain/expense"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

type ExpenseDoc struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Items       string             `bson:"items"`
	Amount      int64              `bson:"amount"`
	PaidDate    time.Time          `bson:"paid_date"`
	PaidBy      string             `bson:"paid_by"`
	Status      string             `bson:"status"`
	DeletedDate *time.Time         `bson:"deleted_date,omitempty"`
}

func NewRepository() (*Repository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get MongoDB URI from environment variable
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	collection := client.Database("expense_tracker").Collection("expenses")
	
	return &Repository{
		client:     client,
		collection: collection,
	}, nil
}

func (r *Repository) Save(exp *expense.Expense) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	doc := ExpenseDoc{
		Items:    exp.Items(),
		Amount:   exp.Amount(),
		PaidDate: exp.PaidDate(),
		PaidBy:   exp.PaidBy(),
		Status:   "active",
	}

	_, err := r.collection.InsertOne(ctx, doc)
	return err
}

func (r *Repository) FindByID(id int) (*expense.Expense, error) {
	// MongoDB uses ObjectID, this is for compatibility
	return nil, nil
}

func (r *Repository) FindAll() ([]*expense.Expense, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var expenses []*expense.Expense
	for cursor.Next(ctx) {
		var doc ExpenseDoc
		if err := cursor.Decode(&doc); err != nil {
			continue
		}
		
		exp := expense.NewExpenseWithDate(doc.Items, doc.Amount, doc.PaidBy, doc.PaidDate)
		expenses = append(expenses, exp)
	}

	return expenses, nil
}

func (r *Repository) FindActiveExpenses() ([]*expense.Expense, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"status": bson.M{"$ne": "deleted"}}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var expenses []*expense.Expense
	for cursor.Next(ctx) {
		var doc ExpenseDoc
		if err := cursor.Decode(&doc); err != nil {
			continue
		}
		
		exp := expense.NewExpenseWithDate(doc.Items, doc.Amount, doc.PaidBy, doc.PaidDate)
		expenses = append(expenses, exp)
	}

	return expenses, nil
}

func (r *Repository) GetAll() ([]map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"status": bson.M{"$ne": "deleted"}}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var expenses []map[string]interface{}
	counter := 1
	for cursor.Next(ctx) {
		var doc ExpenseDoc
		if err := cursor.Decode(&doc); err != nil {
			continue
		}
		
		log.Printf("[MONGO] GetAll - Counter: %d, ObjectID: %s, Items: %s", counter, doc.ID.Hex(), doc.Items)
		expenses = append(expenses, map[string]interface{}{
			"no":       doc.ID.Hex(), // Use ObjectID as string
			"items":    doc.Items,
			"amount":   doc.Amount,
			"paidDate": doc.PaidDate.Format("2006-01-02"),
			"paidBy":   doc.PaidBy,
		})
		counter++
	}

	log.Printf("[MONGO] GetAll - Total expenses returned: %d", len(expenses))
	return expenses, nil
}

func (r *Repository) GetSummaryByPaidBy() (map[string]int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pipeline := []bson.M{
		{"$match": bson.M{"status": bson.M{"$ne": "deleted"}}},
		{"$group": bson.M{
			"_id":   "$paid_by",
			"total": bson.M{"$sum": "$amount"},
		}},
		{"$sort": bson.M{"_id": 1}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	summary := make(map[string]int64)
	for cursor.Next(ctx) {
		var result struct {
			ID    string `bson:"_id"`
			Total int64  `bson:"total"`
		}
		if err := cursor.Decode(&result); err != nil {
			continue
		}
		summary[result.ID] = result.Total
	}

	return summary, nil
}

func (r *Repository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if id == "" {
		return nil
	}

	// Convert string to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("[MONGO] Invalid ObjectID: %s, error: %v", id, err)
		return err
	}

	log.Printf("[MONGO] Soft deleting expense with ObjectID: %s", id)
	filter := bson.M{"_id": objectID}
	now := time.Now()
	update := bson.M{"$set": bson.M{
		"status": "deleted",
		"deleted_date": now,
	}}
	result, err := r.collection.UpdateOne(ctx, filter, update)
	log.Printf("[MONGO] Soft delete result: %+v, error: %v", result, err)
	return err
}

func (r *Repository) GetDeleted() ([]map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"status": "deleted"}
	log.Printf("[MONGO] GetDeleted - Filter: %+v", filter)
	
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		log.Printf("[MONGO] GetDeleted - Find error: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var expenses []map[string]interface{}
	for cursor.Next(ctx) {
		var doc ExpenseDoc
		if err := cursor.Decode(&doc); err != nil {
			log.Printf("[MONGO] GetDeleted - Decode error: %v", err)
			continue
		}
		
		log.Printf("[MONGO] GetDeleted - Found deleted record: ID=%s, Items=%s, Status=%s", doc.ID.Hex(), doc.Items, doc.Status)
		
		deletedDate := "N/A"
		if doc.DeletedDate != nil {
			deletedDate = doc.DeletedDate.Format("2006-01-02")
		}
		
		expenses = append(expenses, map[string]interface{}{
			"id":          doc.ID.Hex(),
			"items":       doc.Items,
			"amount":      doc.Amount,
			"paidDate":    doc.PaidDate.Format("2006-01-02"),
			"paidBy":      doc.PaidBy,
			"deletedDate": deletedDate,
		})
	}

	log.Printf("[MONGO] GetDeleted - Total deleted expenses: %d", len(expenses))
	return expenses, nil
}

func (r *Repository) ClearAll() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Printf("[MONGO] Clearing all expenses from database")
	result, err := r.collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		log.Printf("[MONGO] Clear all error: %v", err)
		return err
	}

	log.Printf("[MONGO] Successfully deleted %d documents", result.DeletedCount)
	return nil
}

func (r *Repository) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return r.client.Disconnect(ctx)
}