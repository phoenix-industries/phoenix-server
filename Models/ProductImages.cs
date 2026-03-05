namespace Phoenix.Models
{
    public class ProductImages
    {
        public int Id { get; set; }
        public Products Product { get; set; }
        public int? ProductId { get; set; } //fk
        public string Image { get; set; }

        public DateTimeOffset CreatedAt { get; set; } = DateTimeOffset.UtcNow;
        public DateTimeOffset UpdatedAt { get; set; } = DateTimeOffset.UtcNow;
        public DateTimeOffset? DeletedAt { get; set; }// nullable
    }
}
