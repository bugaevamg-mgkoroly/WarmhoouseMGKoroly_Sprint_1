from flask import Flask, request, jsonify
import pika
import json

app = Flask(__name__)
devices = {}

def send_mq(msg):
    conn = pika.BlockingConnection(pika.ConnectionParameters('rabbitmq', 5672))
    ch = conn.channel()
    ch.queue_declare(queue='devices')
    ch.basic_publish(exchange='', routing_key='devices', body=json.dumps(msg))
    conn.close()

@app.route('/devices', methods=['POST'])
def register():
    data = request.get_json()
    if not data or 'id' not in data:
        return jsonify({'error': 'ID required'}), 400
    devices[data['id']] = {'status': 'online'}
    send_mq({'id': data['id'], 'action': 'register'})
    return jsonify(devices[data['id']]), 201

@app.route('/devices/<id>', methods=['GET'])
def get(id):
    if id not in devices:
        return jsonify({'error': 'Not found'}), 404
    return jsonify(devices[id])

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)