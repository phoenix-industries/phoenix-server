namespace phoenix.Models
{
    public class User
    {
        public int Id { get; set; }
        public string Name { get; set; }
        public string Email { get; set; }
        public string Phone { get; set; }
        public int Role { get; set; } = 10; // default role is user
        public string City { get; set; }
        public string Governorate { get; set; }
        public DateTimeOffset BirthDate { get; set; }
        public DateTimeOffset CreatedAt { get; set; }= DateTimeOffset.UtcNow;
        public DateTimeOffset UpdatedAt { get; set; }= DateTimeOffset.UtcNow;
        public DateTimeOffset? DeletedAt { get; set; }// nullable

        // navigation properties
        public List<UserBans> ModeratorBans { get; set; }
        public UserBans UserBan { get; set; }
        public List<Products> Products { get; set; }
        public List<ProductReviws> ProductReviws { get; set; }
        public List<Shippings> FromUser { get; set; }
        public List<Shippings> ToUser { get; set; }
        public List<Invoices> Invoices { get; set; }
    }
}
