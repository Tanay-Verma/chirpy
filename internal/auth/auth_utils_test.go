package auth

import "testing"

func TestAuthUtils(t *testing.T) {
	t.Run("Test password hashing", func(t *testing.T) {
		hashedPassword, err := HashPassword("password")
		if err != nil {
			t.Fatal("Error in hasing password:", err)
		}

		err = CheckPasswordHash("password", hashedPassword)
		if err != nil {
			t.Fatal("Error in comparing password and hash:", err)
		}
	})
}
