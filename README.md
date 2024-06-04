# go-api
Teerachai's Go test.

#in main.go line 183 is initial database mock 9 row data every time

#json api collection

{
	"info": {
		"_postman_id": "b073a201-fba6-47ad-9ac6-55d6c77e1193",
		"name": "test api",
		"schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
	},
	"item": [
		{
			"name": "Health check customer",
			"request": {
				"method": "GET",
				"header": [],
				"url": {
					"raw": "http://localhost:8080/",
					"protocol": "http",
					"host": [
						"localhost"
					],
					"port": "8080",
					"path": [
						""
					]
				}
			},
			"response": []
		},
		{
			"name": "Get all customer",
			"request": {
				"method": "GET",
				"header": [],
				"url": {
					"raw": "http://localhost:8080/customers",
					"protocol": "http",
					"host": [
						"localhost"
					],
					"port": "8080",
					"path": [
						"customers"
					]
				}
			},
			"response": []
		},
		{
			"name": "create customer",
			"request": {
				"method": "POST",
				"header": [],
				"body": {
					"mode": "raw",
					"raw": "{\n  \"name\":\"Wick\",\n  \"age\":30\n}",
					"options": {
						"raw": {
							"language": "json"
						}
					}
				},
				"url": {
					"raw": "http://localhost:8080/customers",
					"protocol": "http",
					"host": [
						"localhost"
					],
					"port": "8080",
					"path": [
						"customers"
					]
				}
			},
			"response": []
		},
		{
			"name": "get customer by id",
			"request": {
				"method": "GET",
				"header": [],
				"url": {
					"raw": "http://localhost:8080/customer?id=1",
					"protocol": "http",
					"host": [
						"localhost"
					],
					"port": "8080",
					"path": [
						"customer"
					],
					"query": [
						{
							"key": "id",
							"value": "1"
						}
					]
				}
			},
			"response": []
		},
		{
			"name": "delete customer",
			"request": {
				"method": "DELETE",
				"header": [],
				"url": {
					"raw": "http://localhost:8080/customer?id=2",
					"protocol": "http",
					"host": [
						"localhost"
					],
					"port": "8080",
					"path": [
						"customer"
					],
					"query": [
						{
							"key": "id",
							"value": "2"
						}
					]
				}
			},
			"response": []
		},
		{
			"name": "update customer",
			"request": {
				"method": "PUT",
				"header": [],
				"body": {
					"mode": "raw",
					"raw": "{\n  \"name\":\"Wick\",\n  \"age\":30\n}",
					"options": {
						"raw": {
							"language": "json"
						}
					}
				},
				"url": {
					"raw": "http://localhost:8080/customer?id=20",
					"protocol": "http",
					"host": [
						"localhost"
					],
					"port": "8080",
					"path": [
						"customer"
					],
					"query": [
						{
							"key": "id",
							"value": "20"
						}
					]
				}
			},
			"response": []
		}
	]
}