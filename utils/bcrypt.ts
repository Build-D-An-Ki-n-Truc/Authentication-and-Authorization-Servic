import bcrypt from 'bcrypt'
// This could be other hashing algorithm depend on USER DATABASE

const number = 10

export const hashPassword = async (password: string) => {
    const salt = await bcrypt.genSalt(number)
    return await bcrypt.hash(password, salt)
}

export const comparePassword = async (password: string, hashedPassword: string) => {
    return await bcrypt.compare(password, hashedPassword)
}