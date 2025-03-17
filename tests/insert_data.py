#!/usr/bin/env python3

from pymilvus import connections, Collection, FieldSchema, CollectionSchema, DataType, utility
import time

# 配置
HOST = "localhost"
PORT = "19530"
COLLECTION_NAME = "test_collection"
DIM = 128  # 向量维度
NUM_ENTITIES = 1000  # 插入 1000 条

def create_collection():
    """创建集合"""
    fields = [
        FieldSchema(name="id", dtype=DataType.INT64, is_primary=True, auto_id=True),
        FieldSchema(name="vector", dtype=DataType.FLOAT_VECTOR, dim=DIM)
    ]
    schema = CollectionSchema(fields=fields, description="Test collection")
    collection = Collection(COLLECTION_NAME, schema)
    return collection

def insert_data(collection):
    """插入 1000 条固定数据"""
    vectors = []
    for i in range(NUM_ENTITIES):
        vector = [float(j + 1) for j in range(DIM)]  # [1.0, 2.0, ..., 128.0]
        vectors.append(vector)
    collection.insert([vectors])
    collection.flush()  # 强制写入磁盘
    print(f"Inserted {NUM_ENTITIES} entities into {COLLECTION_NAME}")

def build_index(collection):
    """建立索引"""
    index_params = {
        "metric_type": "L2",
        "index_type": "IVF_FLAT",
        "params": {"nlist": 128}
    }
    collection.create_index("vector", index_params)
    print("Index built successfully")

def main():
    # 连接 Milvus
    connections.connect(host=HOST, port=PORT)
    print("Connected to Milvus")

    # 删除旧集合
    if utility.has_collection(COLLECTION_NAME):
        utility.drop_collection(COLLECTION_NAME)
        print(f"Dropped existing collection: {COLLECTION_NAME}")

    # 创建集合
    collection = create_collection()

    # 插入数据
    insert_data(collection)

    # 建索引
    build_index(collection)

    # 加载集合并检查数量
    collection.load()
    time.sleep(2)  # 等待加载完成
    total = collection.num_entities
    print(f"Total entities: {total}")
    if total != NUM_ENTITIES:
        print(f"Warning: Expected {NUM_ENTITIES}, got {total}")

    connections.disconnect("default")

if __name__ == "__main__":
    main()