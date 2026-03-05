using Azure;

namespace Phoenix.Models
{
    public class ProductTag
    {
        public int ProductId { get; set; }
        public Products Product { get; set; }

        public int TagId { get; set; }
        public Tags Tag { get; set; }
    }
}
