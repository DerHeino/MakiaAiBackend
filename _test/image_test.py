import cv2
import json
import numpy as np
import os
import requests
import sys
import time
from dataclasses import dataclass
from multiprocessing import Process

os.system('color')

@dataclass
class INIT_INFO:
	DOMAIN: str
	HEADER: 'dict'
	DEVICE: str
	TOKEN: str

DOMAIN: str = "https://healthapi-u6xn.onrender.com"

HEADER = {
	"User-Agent":"python-test-routine",
	"Accept-Encoding":"gzip,deflate,br"
}


def loop_images(d: INIT_INFO):

	dir = os.path.join(os.getcwd(), 'images')
	domain = d.DOMAIN + '/' + 'device' + '/' + d.DEVICE + '/' + 'image'

	print('starting upload loop')
	counter = 1
	while True:
		for filename in sorted(os.listdir(dir)):
			f = os.path.join(dir, filename)
			if os.path.isfile(f):
				requests.post(url=domain, files={"image": open(f, 'rb')})
				print(f'image uploaded; counter: {counter}')
				counter += 1
				time.sleep(3)


def get_images(d: INIT_INFO, show: bool):

	domain = d.DOMAIN + '/' + 'device' + '/' + d.DEVICE + '/' + 'image'
	header = d.HEADER.copy()
	header["Authorization"] = f"Bearer {d.TOKEN}"

	print('starting download loop')
	counter = 1
	time.sleep(1)
	while True:
		r : requests.Response = requests.get(url=domain, headers=header)
		if 'image' in r.headers['content-type']:
			if show:
				image = np.asarray(bytearray(r.content), dtype="uint8")
				image = cv2.imdecode(image, cv2.IMREAD_COLOR)
				cv2.imshow('imageloop', image)
				cv2.waitKey(20)
			print(f"image received: {counter}")
			counter += 1
		else:
			print(f"error: not image; image at {counter}")
		time.sleep(2)


def init_models() -> 'tuple[str, str]':
	print('start init')

	f = open('image_init.json', encoding='utf-8')
	init_dict: 'dict' = json.load(f)
	f.close()

	route_login = DOMAIN + '/' + 'login'
	response: requests.Response = requests.post(url=route_login, headers=HEADER, json=init_dict["user"])
	token = response.content.decode('utf8')

	header = HEADER.copy()
	header["Authorization"] = f"Bearer {token}"

	route_project = DOMAIN + '/' + 'project'
	requests.post(url=route_project, headers=header, json=init_dict["project"])

	route_location = DOMAIN + '/' + 'location'
	requests.post(url=route_location, headers=header, json=init_dict["location"])

	route_device = DOMAIN + '/' + 'device'
	requests.post(url=route_device, headers=header, json=init_dict["device"])
	device = init_dict["device"]["_id"]

	return device, token


def args() -> 'tuple[int, bool]':
	t, d = 1, False
	for _, arg in enumerate(sys.argv):
		if arg.isdigit():
			t = int(arg)
			continue
		if arg == 'display':
			d = True
			continue
		if arg == ('True' or 'true'):
			d = True
			continue

	return t, d


def main():
	min, display = args()

	os.chdir(os.path.dirname(os.path.realpath(__file__)))

	device, token = init_models()
	info = INIT_INFO(DOMAIN=DOMAIN, HEADER=HEADER, DEVICE=device, TOKEN=token)

	p1 = Process(target=loop_images, args=[info])
	p1.start()

	#p2 = Process(target=get_images, args=[info, display])
	#p2.start()

	time.sleep(60*min)
	p1.terminate()
	#p2.terminate()

	exit(0)


if __name__ == "__main__":
	main()
