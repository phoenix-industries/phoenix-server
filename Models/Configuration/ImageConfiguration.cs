using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;

namespace phoenix.Models.Configuration
{
    public class ImageConfiguration : IEntityTypeConfiguration<ProductImages>
    {
        public void Configure(EntityTypeBuilder<ProductImages> builder)
        {
            builder.ToTable("ProductImages");
            builder.HasKey(e => e.Id);

            builder.Property(e => e.Image).IsRequired().HasMaxLength(255);
            builder.HasOne(pi => pi.Product)
                   .WithMany(p => p.ProductImages)
                   .HasForeignKey(pi => pi.ProductId)
                   .OnDelete(DeleteBehavior.SetNull);
        }
    }
}
