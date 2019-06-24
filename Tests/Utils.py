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

def generate_region():
    region = dict()
    region['latitude'] = (random.randint(0, 1000)/100)
    region['longitude'] = (random.randint(0, 1000)/100)
    return region

def generate_step():
    step = dict()
    step['stepId'] = (random.randint(0, 100))
    print(step['stepId'])
    step['latlong'] = generate_region()
    return step

def generate_trip():
    trip = dict()
    trip['tripId'] = (random.randint(0, 100))
    trip['steps'] = []
    trip['steps'].append(generate_step())
    trip['steps'].append(generate_step())
    trip['steps'].append(generate_step())
    trip['public'] = True
    return trip