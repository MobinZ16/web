# /// script
# requires-python = ">=3.11"
# dependencies = [
#     "cryptography>=50.0.2",
#     "pyjwt>=2.15.1",
# ]
# ///
import jwt

# همان کلید مخفی که در docker-compose.yml برای پروکسی تعریف کردیم
secret = "IWnqSYO5zVcR3WBVupuxPsL4T1JoRCleZnQ8J3q6uCN"

# تولید توکن معتبر
token = jwt.encode({"sub": "student_test"}, secret, algorithm="HS256")

print("--- TOKEN START ---")
print(token)
print("--- TOKEN END ---")