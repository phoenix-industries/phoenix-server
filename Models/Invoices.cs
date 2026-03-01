namespace phoenix.Models
{
    public enum InvoiceStatus
    {
        Pending,
        Delivered,
        Cancelled
    }
    public class Invoices
    {
        public int Id { get; set; }
        public int UserId { get; set; } //fk
        public User User { get; set; }
        public int productId { get; set; } //fk
        public Products Product { get; set; }
        public InvoiceStatus Status { get; set; } = InvoiceStatus.Pending;
        public string Note { get; set; }
        public DateTimeOffset CreatedAt { get; set; } = DateTimeOffset.UtcNow;
        public DateTimeOffset UpdatedAt { get; set; } = DateTimeOffset.UtcNow;
        public DateTimeOffset? DeletedAt { get; set; } // nullable
    }
}
