namespace Phoenix.Models
{
    public class Tags
    {
        public int Id { get; set; }
        public string Tag { get; set; }
        public DateTimeOffset CreatedAt { get; set; } = DateTimeOffset.UtcNow;
        public DateTimeOffset UpdatedAt { get; set; } = DateTimeOffset.UtcNow;
        public DateTimeOffset? DeletedAt { get; set; }// nullable


        //navigation properties
        public List<ProductTag> ProductTags { get; set; }
    }
}
