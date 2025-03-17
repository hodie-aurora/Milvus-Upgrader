#!/usr/bin/env python3

from pymilvus import connections, Collection, utility
import sys

# 配置
HOST = "localhost"
PORT = "19530"
COLLECTION_NAME = "test_collection"

def clear_data(collection):
    """清空集合中的所有数据，保留集合结构"""
    collection.load()
    total = collection.num_entities
    if total == 0:
        print(f"No data to clear in {COLLECTION_NAME}")
        return
    
    # 删除所有数据
    expr = "id >= 0"  # 删除所有 ID
    collection.delete(expr)
    collection.flush()
    print(f"Cleared {total} entities from {COLLECTION_NAME}")

def drop_collection():
    """删除整个集合"""
    if utility.has_collection(COLLECTION_NAME):
        utility.drop_collection(COLLECTION_NAME)
        print(f"Dropped collection: {COLLECTION_NAME}")
    else:
        print(f"Collection {COLLECTION_NAME} does not exist")

def main():
    # 连接 Milvus
    connections.connect(host=HOST, port=PORT)
    print("Connected to Milvus")

    if not utility.has_collection(COLLECTION_NAME):
        print(f"Collection {COLLECTION_NAME} not found")
        connections.disconnect("default")
        return

    # 提供选项
    print("Choose an action:")
    print("1. Clear all data (keep collection)")
    print("2. Drop entire collection")
    choice = input("Enter 1 or 2: ")

    collection = Collection(COLLECTION_NAME)

    if choice == "1":
        clear_data(collection)
        total = collection.num_entities
        print(f"After clearing, total entities: {total}")
    elif choice == "2":
        drop_collection()
    else:
        print("Invalid choice")

    connections.disconnect("default")

if __name__ == "__main__":
    main()