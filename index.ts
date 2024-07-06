import express, { Express, Request, Response , Application } from 'express';
import {config} from './config/config';

// creating express app
const app: Application = express();

app.use(express.json())

const port = config.port || 8000;

app.get('/', (req: Request, res: Response) => {
  res.send('Welcome to Express & TypeScript Server');
});

app.listen(port, () => {
  console.log(`Server is Fire at http://localhost:${port}`);
});