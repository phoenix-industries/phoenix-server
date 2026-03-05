using Microsoft.EntityFrameworkCore;
using Phoenix.Models.Configuration;

namespace Phoenix.Models
{
    public class Context : DbContext
    {
		private readonly IConfiguration _configuration;

        public DbSet<User> Users { get; set; }
        public DbSet<Products> Products { get; set; }
        public DbSet<ProductImages> ProductImages { get; set; }
        public DbSet<ProductCategories> Categories { get; set; }
        public DbSet<ProductTag> productTags { get; set; }
        public DbSet<ProductReviws> ProductReviws { get; set; }
        public DbSet<Tags> Tags { get; set; }
        public DbSet<UserBans> UserBans { get; set; }
        public DbSet<Shippings> Shippings { get; set; }
        public DbSet<Invoices> Invoices { get; set; }

		public Context(DbContextOptions<Context> options, IConfiguration configuration) : base(options) {
			this._configuration = configuration;
		}

        protected override void OnConfiguring(DbContextOptionsBuilder optionsBuilder)
        {
			optionsBuilder.UseNpgsql(this._configuration.GetConnectionString("DefaultConnection"));
        }

        protected override void OnModelCreating(ModelBuilder modelBuilder)
        {
            modelBuilder.ApplyConfiguration(new UserConfiguration());
            modelBuilder.ApplyConfiguration(new ProductConfiguration());
            modelBuilder.ApplyConfiguration(new ImageConfiguration());
            modelBuilder.ApplyConfiguration(new CategoryConfiguration());
            modelBuilder.ApplyConfiguration(new ProductTagsConfiguration());
            modelBuilder.ApplyConfiguration(new ReviewConfiguration());
            modelBuilder.ApplyConfiguration(new UserBansConfiguration());
            modelBuilder.ApplyConfiguration(new ShippingConfiguration());
            modelBuilder.ApplyConfiguration(new InvoicesConfiguration());

            base.OnModelCreating(modelBuilder);
        }
    }
}
