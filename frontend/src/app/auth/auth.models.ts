export interface LoginPayload {
  role: 'Employee' | 'Manager';
  employeeId: string;
  password: string;
}
