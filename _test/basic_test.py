import base64
import datetime
import requests
import os
import time
import json
from typing import Callable

os.system('color')

# request header with information
# header auth after login
# supports json data

DOMAIN = None

# default header contains User-Agent, Accept-Encoding, Accept, Connection
# I'm adding my own User-Agent and Accept-Encoding
# Content-Type should be added automatically and Authorization Token will be added after login
HEADER = {
	"User-Agent":"python-test-routine",
	"Accept-Encoding":"gzip,deflate,br"
}

USER_DICT	: 'dict'
VALID_DICT	: 'dict'
INVALID_DICT: 'dict'

ADMIN_TOKEN	: 'str' = None
USER_TOKEN	: 'str'	= None


SUCCESS = "SUCCESS"
FAIL 	= "FAILURE"
SKIPPED	= "SKIPPED"

RESPONSE_ERR = "error"


def check_JWT(response: str, _) -> bool:

	def pad_b64(s: str) -> str:
		i = len(s) % 4
		for _ in range(i):
			s += '='
		return s

	parts: 'list' = response.split('.')
	if len(parts) != 3:
		return False
	part_one: 'dict' = json.loads(base64.b64decode(pad_b64(parts[0])))
	if ('alg' or 'typ') not in part_one.keys():
		return False
	part_two: 'dict' = json.loads(base64.b64decode(pad_b64(parts[1])))
	if 'username' not in part_two.keys():
		return False
	return True


def check_UUID(response: str, _) -> bool:

	parts: 'list' = response.split('-')
	if len(parts) != 5:
		print(len(parts), parts, 'here')
		return False
	if len(parts[0]) != 8 or len(parts[1]) != 4 or len(parts[2]) != 4 or len(parts[3]) != 4 or len(parts[4]) != 12:
		return False
	return True


def equal(response: str, expected: str) -> bool:
	if response == expected:
		return True


def contains(response: str, expected: str) -> bool:
	if expected in response:
		return True


def check_response(test: str, content: bytes, expected: str, method: Callable = equal) -> bool:
	response = content.decode('utf8')

	if method(response, expected):
		if test: print(f'[ \x1b[0;32;40m{SUCCESS}\x1b[0m ]\t{test}')
		if expected == RESPONSE_ERR: print(f'\t\t{response}')
		return True

	if test: print(f'[ \x1b[0;31;40m{FAIL}\x1b[0m ]\t{test}\n\t\t{response}')
	return False


def check_skipped(*skipped: 'str'):
	for s in skipped:
		print(f'[ \x1b[0;33;40m{SKIPPED}\x1b[0m ]\t{s}')


def convert_time(time_format: str) -> str:

	if time_format == "$isoTimestamp":
		return datetime.datetime.utcnow().replace(microsecond=0).isoformat() + 'Z'
	else:
		unix_time = time.time()
		return str(unix_time).split('.')[0]


def login() -> bool:
	print(f"======== {'LOGIN TESTING':16s} ========")
	global ADMIN_TOKEN
	global USER_TOKEN

	response = make_request(func=requests.post, route="login", jsonObj=USER_DICT["login_admin"])
	if not check_response("login admin", response.content, "", check_JWT):
		return False

	ADMIN_TOKEN = response.content.decode('utf8')

	response = make_request(func=requests.post, route="login", jsonObj=USER_DICT["login_user"])
	if not check_response("login user", response.content, "", check_JWT):
		return False

	USER_TOKEN = response.content.decode('utf8')

	response = make_request(func=requests.post, route="login", jsonObj=USER_DICT["login_invalid"])
	check_response("login invalid body", response.content, RESPONSE_ERR, contains)

	return True


def register(time: bool):
	print(f"======== {'REGISTER TESTING':16s} ========")
	otk: str = None
	delete: bool = False

	response = make_request(func=requests.get, route="register", token=ADMIN_TOKEN)
	if check_response("register get key", response.content, "", check_JWT):
		otk = response.content.decode('utf8')

	response = make_request(func=requests.get, route="register", token=USER_TOKEN)
	check_response("register get key invalid", response.content, RESPONSE_ERR, contains)

	if otk:
		response = make_request(func=requests.post, route="register", token=otk, jsonObj=USER_DICT["post_register"])
		if check_response("register user success", response.content, "", check_JWT):
			delete = True
	else:
		check_skipped("register user success")

	response = make_request(func=requests.post, route="register", jsonObj=USER_DICT["login_invalid"])
	check_response("register no token", response.content, RESPONSE_ERR, contains)

	response = make_request(func=requests.post, route="register", token=USER_TOKEN, jsonObj=USER_DICT["login_invalid"])
	check_response("register invalid token", response.content, RESPONSE_ERR, contains)

	response = make_request(func=requests.get, route="register", token=ADMIN_TOKEN)
	if check_response(None, response.content, "", check_JWT):
		otk = response.content.decode('utf8')

		response = make_request(func=requests.post, route="register", jsonObj=USER_DICT["login_invalid"])
		check_response("register empty body", response.content, RESPONSE_ERR, contains)

		response = make_request(func=requests.post, route="register", jsonObj=USER_DICT["login_invalid"])
		check_response("register invalid body", response.content, RESPONSE_ERR, contains)

		response = make_request(func=requests.post, route="register", jsonObj=USER_DICT["post_register"])
		check_response("register already exists", response.content, RESPONSE_ERR, contains)

		if time:
			time.sleep(100)
			response = make_request(func=requests.post, route="register", jsonObj=USER_DICT["post_register"])
			check_response("register timeout", response.content, RESPONSE_ERR, contains)
		else:
			check_skipped("register timeout")
	else:
		check_skipped("register empty body", "register invalid body", "register already exists", "register timeout")

	# delete user
	if delete:
		pass
		#response = make_request(func=requests.delete, route="login", token=ADMIN_TOKEN, data=USER_DICT["post_register"]["username"])
		#check_response("login_delete", response.content, "", contains)
	else:
		check_skipped("delete user")


def post_valid():

	def get_method(identifier: str) -> Callable:
		if identifier == "equal":
			return equal
		elif identifier == "contains":
			return contains
		elif identifier == "UUID":
			return check_UUID
		else:
			return None

	for key in VALID_DICT:
		model_list = VALID_DICT[key]
		token = ADMIN_TOKEN if key in ('project', 'location') else USER_TOKEN

		print(f"======== POST {key:11s} ========")

		i = 1
		for model in model_list:
			check = model.pop('_return')
			expected = model['_id'] if key != 'ping' else model['id']
			if key in ('ping', 'inventory'):
				for model_key in model:
					if '$' in model[model_key]:
						model[model_key] = convert_time(model[model_key])
			if key == 'device':
				if 'lastPing' in model:
					if 'timestamp' in model['lastPing']:
						model['lastPing']['timestamp'] = convert_time(model['lastPing']['timestamp'])
			method = get_method(check)

			response = make_request(func=requests.post, route=key, token=token, jsonObj=model)
			check_response(f"{key} {i}", response.content, expected, method)

			i += 1


def post_invalid():
	
	for key in INVALID_DICT:
		model_list = INVALID_DICT[key]
		token = ADMIN_TOKEN if key in ('project', 'location') else USER_TOKEN

		print(f"======== ERROR {key:10s} ========")
		
		i = 1
		for model in model_list:
			if key in ('ping', 'inventory'):
				for model_key in model:
					if type(model[model_key]) == str:
						if '$' in model[model_key]:
							model[model_key] = convert_time(model[model_key])
			if key == 'device':
				if 'lastPing' in model:
					if 'timestamp' in model['lastPing']:
						model['lastPing']['timestamp'] = convert_time(model['lastPing']['timestamp'])

			response = make_request(func=requests.post, route=key, token=token, jsonObj=model)
			check_response(f"{key} {i}", response.content, RESPONSE_ERR, contains)

			i += 1


def delete():
	pass
	# inventory: dd5891ae-c966-42be-97ee-811fca8f90ac
	# device: ce0c8bee-9844-4ff8-a842-6d506e8f2b65
	# location: cf46a3c1-1ccb-45c6-bd0b-321271c69b87 (not nested after del ce0c8bee)
	# project: id: fish (get id from get all)

	# device: 9eb5c5dd-21a7-4020-8c67-a24eaea16790 (nested)
	# location: 58a5ec1e-c5bf-4463-9c55-e4735bbf8422 (nested)
	# project: 42df27c3-21e9-4ac8-83b3-24a8176dd825 (nested)


def make_request(func: Callable, route: str, token: str = None, jsonObj: 'dict' = None, data: 'str' = None) -> requests.Response:
	url = DOMAIN + '/' + route
	header = HEADER.copy()
	if token:
		header["Authorization"] = f"Bearer {token}"
	
	return func(url, headers=header, json=jsonObj, data=data)


def open_files():
	global USER_DICT
	global VALID_DICT
	global INVALID_DICT

	f1 = open('user.json', encoding='utf-8')
	USER_DICT = json.load(f1)
	f1.close()

	f2 = open('post_valid.json', encoding='utf-8')
	VALID_DICT = json.load(f2)
	f2.close()

	f3 = open('post_invalid.json', encoding='utf-8')
	INVALID_DICT = json.load(f3)
	f3.close()


def main():
	global DOMAIN

	os.chdir(os.path.dirname(os.path.realpath(__file__)))
	DOMAIN = "http://localhost:10000"

	open_files()
	if login():
		register(time=False)
		post_valid()
		post_invalid()
		# delete()


if __name__ == "__main__":
    main()
