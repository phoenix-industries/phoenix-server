namespace phoenix.Models
{
    public class ProductReviws
    {
        public int Id { get; set; }
        public int? ProductId { get; set; } //fk
        public Products Product { get; set; }
        public int? UserId { get; set; } //fk
        public User User { get; set; }
        public string Comment { get; set; }
        public int Rating { get; set; } 
        public DateTimeOffset CreatedAt { get; set; } = DateTimeOffset.UtcNow;
        public DateTimeOffset UpdatedAt { get; set; } = DateTimeOffset.UtcNow;
        public DateTimeOffset? DeletedAt { get; set; } // nullable
    }
}
