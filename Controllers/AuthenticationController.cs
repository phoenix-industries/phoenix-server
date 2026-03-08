using System.Security.Cryptography;
using System.ComponentModel.DataAnnotations;
using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Authentication.JwtBearer;
using Phoenix.Models;
using Phoenix.Utils.Scrypt;

namespace Phoenix.Server.Controllers
{
    public class UserRegisterData
    {
        [StringLength(255, ErrorMessage = "{0} length must be between {2} and {1}.", MinimumLength = 2)]
        public string Name { get; set; }
        [RegularExpression(@"^((?!\.)[\w\-_.]*[^.])(@\w+)(\.\w+(\.\w+)?[^.\W])$")]
        public string Email { get; set; }
        [StringLength(13, ErrorMessage = "{0} length must be between {2} and {1}.", MinimumLength = 9)]
        public string Phone { get; set; }
        public string City { get; set; }
        public string Governorate { get; set; }
        [StringLength(15, ErrorMessage = "{0} length must be between {2} and {1}.", MinimumLength = 10)]
        public string NationalID { get; set; }
		[RegularExpression(@"^(?=.*\d)(?=.*[A-Z])(?=.*[a-z])(?=.*[^\w\d\s:])([^\s]){8,16}$")]
        public string Password { get; set; }
        public DateTimeOffset BirthDate { get; set; }
    }

    [ApiController]
    [Route("api/[controller]")]
    public class AuthenticationController : Controller
    {
		private readonly Context _context;

		public AuthenticationController(Context context) {
			this._context = context;
		}

        [HttpPost("register")]
        public async Task<IActionResult> Register(UserRegisterData data)
        {
            var res = _context.Users.Where(u => u.NationalID == data.NationalID || u.Email == data.Email || u.Phone == data.Phone);
            if (res != null)
            {
				Console.WriteLine(res.Count());
                return BadRequest("User already exists");
            }

			var hash = ScryptPasswordHasher.Hash(data.Password);

            var user = new User
            {
                Name = data.Name,
                Email = data.Email,
                Phone = data.Phone,
                City = data.City,
                Governorate = data.Governorate,
                NationalID = data.NationalID,
                BirthDate = data.BirthDate,
				Password = hash,
            };

            return Ok(user);
        }
    }
}
