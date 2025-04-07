import { User } from "../types";

export function greetUser(user: User): string {
  return `Hello, ${user.name}!`;
}