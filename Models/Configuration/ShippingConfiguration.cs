using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;

namespace phoenix.Models.Configuration
{
    public class ShippingConfiguration : IEntityTypeConfiguration<Shippings>
    {
        public void Configure(EntityTypeBuilder<Shippings> builder)
        {
            builder.ToTable("Shippings");
            builder.HasKey(s => s.Id);
            builder.Property(s => s.Location).IsRequired().HasMaxLength(1000);

                builder.HasOne(s => s.From)
                    .WithMany(s => s.FromUser)
                    .HasForeignKey(s => s.FromId)
                    .OnDelete(DeleteBehavior.Restrict);

                builder.HasOne(s => s.To)
                    .WithMany(s => s.ToUser)
                    .HasForeignKey(s => s.ToId)
                    .OnDelete(DeleteBehavior.Restrict);
        }
    }
}
