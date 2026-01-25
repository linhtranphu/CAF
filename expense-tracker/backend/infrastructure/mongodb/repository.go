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
	settings   *mongo.Collection
	users      *mongo.Collection
}

type ExpenseDoc struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	Items           string             `bson:"items"`
	Amount          int64              `bson:"amount"`
	Quantity        string             `bson:"quantity,omitempty"`
	Unit            string             `bson:"unit,omitempty"`
	BaseQuantity    string             `bson:"base_quantity,omitempty"`
	BaseUnit        string             `bson:"base_unit,omitempty"`
	OriginalMessage string             `bson:"original_message,omitempty"`
	PaidDate        time.Time          `bson:"paid_date"`
	PaidBy          string             `bson:"paid_by"`
	Status          string             `bson:"status"`
	DeletedDate     *time.Time         `bson:"deleted_date,omitempty"`
}

type UserDoc struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Username  string             `bson:"username"`
	Password  string             `bson:"password"`
	Role      string             `bson:"role"`
	CreatedAt time.Time          `bson:"created_at"`
}

func NewRepository() (*Repository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

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
	settings := client.Database("expense_tracker").Collection("settings")
	users := client.Database("expense_tracker").Collection("users")
	
	return &Repository{
		client:     client,
		collection: collection,
		settings:   settings,
		users:      users,
	}, nil
}

func (r *Repository) Save(exp *expense.Expense) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	doc := ExpenseDoc{
		Items:           exp.Items(),
		Amount:          exp.Amount(),
		Quantity:        exp.Quantity(),
		Unit:            exp.Unit(),
		BaseQuantity:    exp.BaseQuantity(),
		BaseUnit:        exp.BaseUnit(),
		OriginalMessage: exp.OriginalMessage(),
		PaidDate:        exp.PaidDate(),
		PaidBy:          exp.PaidBy(),
		Status:          "active",
	}

	log.Printf("[MONGO] Saving expense: Items=%s, Quantity=%s, Unit=%s, BaseQuantity=%s, BaseUnit=%s", 
		doc.Items, doc.Quantity, doc.Unit, doc.BaseQuantity, doc.BaseUnit)

	_, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		log.Printf("[MONGO] Save error: %v", err)
	}
	return err
}

func (r *Repository) FindByID(id int) (*expense.Expense, error) {
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
			"no":              doc.ID.Hex(),
			"items":           doc.Items,
			"amount":          doc.Amount,
			"quantity":        doc.Quantity,
			"unit":            doc.Unit,
			"baseQuantity":    doc.BaseQuantity,
			"baseUnit":        doc.BaseUnit,
			"originalMessage": doc.OriginalMessage,
			"paidDate":        doc.PaidDate.Format("2006-01-02"),
			"paidBy":          doc.PaidBy,
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
			"id":              doc.ID.Hex(),
			"items":           doc.Items,
			"amount":          doc.Amount,
			"quantity":        doc.Quantity,
			"unit":            doc.Unit,
			"baseQuantity":    doc.BaseQuantity,
			"baseUnit":        doc.BaseUnit,
			"originalMessage": doc.OriginalMessage,
			"paidDate":        doc.PaidDate.Format("2006-01-02"),
			"paidBy":          doc.PaidBy,
			"deletedDate":     deletedDate,
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

func (r *Repository) SaveAPIKey(apiKey string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"key": "gemini_api_key"}
	update := bson.M{"$set": bson.M{
		"key":        "gemini_api_key",
		"value":      apiKey,
		"updated_at": time.Now(),
	}}
	opts := options.Update().SetUpsert(true)
	
	_, err := r.settings.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Printf("[MONGO] Save API key error: %v", err)
		return err
	}
	
	log.Printf("[MONGO] API key saved successfully")
	return nil
}

func (r *Repository) GetAPIKey() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result struct {
		Value string `bson:"value"`
	}
	
	filter := bson.M{"key": "gemini_api_key"}
	err := r.settings.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", nil
		}
		log.Printf("[MONGO] Get API key error: %v", err)
		return "", err
	}
	
	return result.Value, nil
}

func (r *Repository) CreateUser(username, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if user exists
	var existing UserDoc
	err := r.users.FindOne(ctx, bson.M{"username": username}).Decode(&existing)
	if err == nil {
		return mongo.ErrNoDocuments // User exists
	}

	// Determine role
	role := "supervisor"
	if username == "admin" {
		role = "admin"
	}

	doc := UserDoc{
		Username:  username,
		Password:  password,
		Role:      role,
		CreatedAt: time.Now(),
	}

	_, err = r.users.InsertOne(ctx, doc)
	return err
}

func (r *Repository) GetUser(username string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user UserDoc
	err := r.users.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return "", err
	}
	return user.Password, nil
}

func (r *Repository) InitDefaultUsers() error {
	defaultUsers := map[string]string{
		"admin": "admin123",
		"linh":  "linh123",
		"toan":  "toan123",
	}

	for username, password := range defaultUsers {
		err := r.CreateUser(username, password)
		if err != nil && err != mongo.ErrNoDocuments {
			log.Printf("[MONGO] Failed to create user %s: %v", username, err)
		}
	}
	return nil
}
