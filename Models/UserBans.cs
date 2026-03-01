namespace phoenix.Models
{
    public class UserBans
    {
        public int Id { get; set; }
        public int? UserId { get; set; } //fk
        public User User { get; set; }
        public int? ModeratorId { get; set; } //fk
        public User Moderator { get; set; }
        public string Reason { get; set; }
        public DateTimeOffset CreatedAt { get; set; } = DateTimeOffset.UtcNow;
        public DateTimeOffset UpdatedAt { get; set; } = DateTimeOffset.UtcNow;
        public DateTimeOffset? DeletedAt { get; set; }// nullable
    }
}
