namespace Phoenix.Models
{
    public class ProductCategories
    {
        public int Id { get; set; }
        public required string Category { get; set; }
        public DateTimeOffset CreatedAt { get; set; } = DateTimeOffset.UtcNow;
        public DateTimeOffset UpdatedAt { get; set; } = DateTimeOffset.UtcNow;
        public DateTimeOffset? DeletedAt { get; set; }// nullable


        //navigation properties
        public List<Products> Products { get; set; }

    }
}
