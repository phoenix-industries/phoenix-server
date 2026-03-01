using Microsoft.EntityFrameworkCore.Metadata.Builders;
using Microsoft.EntityFrameworkCore;
namespace phoenix.Models.Configuration
{
    public class CategoryConfiguration : IEntityTypeConfiguration<ProductCategories>
    {
        public void Configure(EntityTypeBuilder<ProductCategories> builder)
        {
            builder.ToTable("ProductCategories");
            builder.HasKey(e => e.Id);
            builder.Property(e => e.Category).IsRequired().HasMaxLength(100);
        }       
    }
}
