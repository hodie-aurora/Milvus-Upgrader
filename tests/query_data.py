#!/usr/bin/env python3

from pymilvus import connections, Collection, utility
import time

# 配置
HOST = "localhost"
PORT = "19530"
COLLECTION_NAME = "test_collection"

def query_data():
    """查询并验证数据"""
    connections.connect(host=HOST, port=PORT)
    print("Connected to Milvus")

    if not utility.has_collection(COLLECTION_NAME):
        print(f"Collection {COLLECTION_NAME} not found")
        return

    collection = Collection(COLLECTION_NAME)
    collection.load()
    time.sleep(2)  # 等待加载完成

    # 获取总数
    total = collection.num_entities
    print(f"Total entities in {COLLECTION_NAME}: {total}")

    # 查询前 5 条数据
    results = collection.query(expr="id >= 0", limit=5, output_fields=["id", "vector"])
    print("Sample data (first 5):")
    for result in results:
        print(f"ID: {result['id']}, Vector (first 5 dims): {result['vector'][:5]}")

    # 验证总数
    if total == 1000:
        print("Data verification passed: 1000 entities found")
    else:
        print(f"Data verification failed: Expected 1000, got {total}")

    connections.disconnect("default")

if __name__ == "__main__":
    query_data()