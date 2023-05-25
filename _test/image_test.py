import json
import os
import requests
from multiprocessing import Process

os.system('color')

DOMAIN = None

HEADER = {
	"User-Agent":"python-test-routine",
	"Accept-Encoding":"gzip,deflate,br"
}


def loop_images():
	pass


def get_images():
	pass


def init_models():
	f = open('image_init.json', encoding='utf-8')
	init_dict: 'dict' = json.load(f)
	f.close()

	route_login = DOMAIN + '/' + 'login'
	response: requests.Response = requests.post(route=route_login, headers=HEADER, jsonObj=init_dict["user"])
	token = response.content.decode('utf8')

	header = HEADER.copy()
	header["Authorization"] = f"Bearer {token}"

	route_project = DOMAIN + '/' + 'project'
	requests.post(route=route_project, headers=header, json=init_dict["project"])

	route_location = DOMAIN + '/' + 'location'
	requests.post(route=route_location, headers=header, json=init_dict["location"])

	route_device = DOMAIN + '/' + 'device'
	requests.post(route=route_device, headers=header, json=init_dict["device"])


def main():
	global DOMAIN

	os.chdir(os.path.dirname(os.path.realpath(__file__)))
	DOMAIN = "http://localhost:10000"

	init_models()

	p = Process(target=loop_images)
	p.start()

	get_images()

	# command line argument number in minutes (default 1 min)


if __name__ == "__main__":
    main()
