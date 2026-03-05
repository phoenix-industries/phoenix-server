using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;

namespace Phoenix.Models.Configuration
{
    public class UserBansConfiguration : IEntityTypeConfiguration<UserBans>
    {
        public void Configure(EntityTypeBuilder<UserBans> builder)
        {
            builder.ToTable("UserBans");
            builder.HasKey(ub => ub.Id);


            // Relationships

            //relation with moderator
            builder.HasOne(u => u.Moderator)
                .WithMany(u => u.ModeratorBans)
                .HasForeignKey(u => u.ModeratorId)
                .OnDelete(DeleteBehavior.Restrict);

            //relation with user
            builder.HasOne(u => u.User)
                .WithOne(u => u.UserBan)
                .HasForeignKey<UserBans>(u => u.UserId)
                .OnDelete(DeleteBehavior.Restrict);

        }
    }
}
