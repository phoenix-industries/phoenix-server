using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;

namespace phoenix.Models.Configuration
{
    public class InvoicesConfiguration : IEntityTypeConfiguration<Invoices>
    {
        public void Configure(EntityTypeBuilder<Invoices> builder)
        {
            builder.ToTable("Invoices");
            builder.HasKey(i => i.Id);

            builder.HasOne(i => i.User)
                    .WithMany(u => u.Invoices)
                    .HasForeignKey(i => i.UserId)
                    .OnDelete(DeleteBehavior.Restrict);

            builder.HasOne(i => i.Product)
                    .WithMany(p => p.Invoices)
                    .HasForeignKey(i => i.productId)
                    .OnDelete(DeleteBehavior.Restrict);
        }
    }
}
