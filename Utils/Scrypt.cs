using System.Security.Cryptography;
using CryptSharp;

namespace Phoenix.Utils.Scrypt
{
    public static class ScryptPasswordHasher
    {
        public static string Hash(string password)
        {
            if (string.IsNullOrWhiteSpace(password))
            {
                throw new ArgumentException("Password cannot be empty.");
            }
            return Crypter.Blowfish.Crypt(password);
        }

		/*
        public static bool Verify(string password, string hashedPassword)
        {
            if (string.IsNullOrWhiteSpace(password) || string.IsNullOrWhiteSpace(hashedPassword))
            {
                return false;
            }
            return Encoder.Compare(password, hashedPassword);
        }
		*/
    }
}
