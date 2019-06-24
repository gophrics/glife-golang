import string
import random
def generate_name():
    letters = string.ascii_lowercase
    return ''.join(random.choice(letters) for i in range(5)) + ' ' + ''.join(random.choice(letters) for i in range(7))

def generate_username():
    letters = string.ascii_lowercase
    return ''.join(random.choice(letters) for i in range(9))

def generate_email():
    letters = string.ascii_lowercase
    return ''.join(random.choice(letters) for i in range(6)) + '@' + ''.join(random.choice(letters) for i in range(6)) + '.com'

def generate_password():
    letters = string.ascii_lowercase
    return ''.join(random.choice(letters) for i in range(6))

def generate_phone():
    letters = string.digits
    return ''.join(random.choice(letters) for i in range(6))

def generate_country():
    return "India"