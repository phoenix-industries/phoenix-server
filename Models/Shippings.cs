using System.Runtime.CompilerServices;

namespace Phoenix.Models
{
    public class Shippings
    {
        public int Id { get; set; }
        public int FromId { get; set; } //fk
        public User From { get; set; }
        public int ToId { get; set; } //fk
        public User To { get; set; }
        public float Fee { get; set; }
        public String Note { get; set; }
        public string Location { get; set; }
        public DateTimeOffset CreatedAt { get; set; } = DateTimeOffset.UtcNow;
        public DateTimeOffset UpdatedAt { get; set; } = DateTimeOffset.UtcNow;
        public DateTimeOffset? DeletedAt { get; set; } // nullable
    }
}
