package models

type Book struct {
	ID      string   `json:"id" bson:"_id"`
	Title   string   `json:"title"`
	Author  string   `json:"author"`
	Amount  int      `json:"amount"`
	OwnedBy []string `json:"owned_by" bson:"owned_by"`
}

type User struct {
	ID              string   `json:"id" bson:"-"`
	Username        string   `json:"username"`
	Password        string   `json:"-"`
	Role            string   `json:"role"`
	BorrowedBookIDs []string `json:"borrowed_book_ids" bson:"borrowed_book_ids"`
}
