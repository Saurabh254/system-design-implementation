import schemas


def get_user_feed(user_id: int, pagination: schemas.Pagination, db):
    
    return {"user_id": user_id, "feed": ["post1", "post2", "post3"]}