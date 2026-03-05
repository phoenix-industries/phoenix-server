namespace Phoenix.Models
{
    public class Products
    {
        public int Id { get; set; }
        public string Name { get; set; }
        public User User { get; set; }
        public int? UserId { get; set; } //fk
        public ProductCategories ProductCategory { get; set; }
        public int? ProductCategoryId { get; set; } //fk
        public DateTimeOffset CreatedAt { get; set; } = DateTimeOffset.UtcNow;
        public DateTimeOffset UpdatedAt { get; set; } = DateTimeOffset.UtcNow;
        public DateTimeOffset? DeletedAt { get; set; }// nullable

        //navigation properties
        public List<ProductImages> ProductImages { get; set; }
        public List<ProductTag> ProductTags { get; set; }
        public List<ProductReviws> ProductReviws { get; set; }
        public List<Invoices> Invoices { get; set; }
    }
}
