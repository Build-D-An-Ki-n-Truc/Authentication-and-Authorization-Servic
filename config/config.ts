import dotenv from "dotenv"

dotenv.config()

export const config = {
    port: process.env.PORT,
    jwtSecret: process.env.JWT_SECRET,
    dbUrl: process.env.DB_URL,
}