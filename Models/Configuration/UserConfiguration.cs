using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;

namespace Phoenix.Models.Configuration
{
    public class UserConfiguration : IEntityTypeConfiguration<User>
    {
        public void Configure(EntityTypeBuilder<User> builder)
        {
            builder.ToTable("Users");
            builder.HasKey(u => u.Id);
            builder.Property(u => u.Name).IsRequired().HasMaxLength(100);
            builder.Property(u => u.Email).IsRequired().HasMaxLength(100);
            builder.Property(u => u.Phone).HasMaxLength(20);

            //indexes
            builder.HasIndex(u => u.Email).IsUnique().HasDatabaseName("users_email_idx");
            builder.HasIndex(u => u.Phone).IsUnique().HasDatabaseName("users_phone_idx");
        }
    }
}
