from fastapi import FastAPI, Depends 
import service 
from database import get_db
from sqlalchemy.orm import Session
import schemas
app = FastAPI()



@app.get("/healthz")
async def health_check():
    return {"status": "ok"}



@app.get("/user/{user_id}/feed")
async def get_user_feed(user_id: int, pagination: schemas.Pagination = Depends(), db: Session = Depends(get_db)):

    return service.get_user_feed(user_id, pagination, db)