using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;

namespace phoenix.Models.Configuration
{
    public class ReviewConfiguration : IEntityTypeConfiguration<ProductReviws>
    {
        public void Configure(EntityTypeBuilder<ProductReviws> builder)
        {
            builder.ToTable("ProductReviws");
            builder.HasKey(r =>r.Id);
            builder.Property(r => r.Comment).IsRequired().HasMaxLength(1000);

            builder.HasOne(r => r.User)
                .WithMany(r => r.ProductReviws)
                .HasForeignKey(r => r.UserId)
                .OnDelete(DeleteBehavior.SetNull);

            builder.HasOne(r => r.Product)
                .WithMany(r => r.ProductReviws)
                .HasForeignKey(r => r.ProductId)
                .OnDelete(DeleteBehavior.SetNull);
        }
    }
}
