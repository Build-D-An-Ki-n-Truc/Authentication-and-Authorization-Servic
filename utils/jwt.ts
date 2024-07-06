import jwt from 'jsonwebtoken'
import { config } from '../config/config'

const secret: string = config.jwtSecret || ""

export const generateToken = (payload: any) => {
    return jwt.sign(
        payload, 
        secret, 
        {expiresIn: '10m'},
    )
}

export const verifyToken = (token: string) => {
    return jwt.verify(
        token,
        secret,
    )
}
