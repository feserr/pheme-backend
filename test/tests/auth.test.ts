import { describe, it, expect } from '@jest/globals';
import request from 'supertest';

const baseUrl = 'http://127.0.0.1:8000';

async function GetCookie() {
  const response = await request(baseUrl)
    .post('/api/v1/auth/login')
    .send({ email: 'test@auth.com', password: 'test' });

  return response.get('Set-Cookie')[0];
}

describe('Register endpoint', () => {
  it('register new user missing name', async () => {
    const response = await request(baseUrl)
      .post('/api/v1/auth/register')
      .send({ email: 'test@auth.com', password: 'test' });

    expect(response.statusCode).toBe(400);
  });

  it('register new user missing email', async () => {
    const response = await request(baseUrl)
      .post('/api/v1/auth/register')
      .send({ username: 'test', password: 'test' });

    expect(response.statusCode).toBe(400);
  });

  it('register new user missing password', async () => {
    const response = await request(baseUrl)
      .post('/api/v1/auth/register')
      .send({ username: 'test', email: 'test@auth.com' });

    expect(response.statusCode).toBe(400);
  });

  it('register new user', async () => {
    const response = await request(baseUrl)
      .post('/api/v1/auth/register')
      .send({ username: 'test', email: 'test@auth.com', password: 'test' });

    expect(response.statusCode).toBe(200);
    expect(response.headers['content-type']).toContain('application/json');
    expect(response.body).toHaveProperty('id');
    expect(response.body).toHaveProperty('version');
    expect(response.body).toHaveProperty('username');
    expect(response.body).toHaveProperty('email');
    expect(response.body).toHaveProperty('createdAt');
  });
});

describe('Register endpoint with data', () => {
  it('register duplicate user', async () => {
    const response = await request(baseUrl)
      .post('/api/v1/auth/register')
      .send({ username: 'test', email: 'test@auth.com', password: 'test' });

    expect(response.statusCode).toBe(400);
  });
});

describe('Login endpoint', () => {
  it('login with existing user', async () => {
    const response = await request(baseUrl)
      .post('/api/v1/auth/login')
      .send({ email: 'test@auth.com', password: 'test' });

    expect(response.statusCode).toBe(200);
    expect(response.get('Set-Cookie')).toBeDefined();
  });

  it('get logged user info', async () => {
    const cookie = await GetCookie();

    expect(cookie).toBeDefined();

    const response = await request(baseUrl)
      .get('/api/v1/auth/user')
      .set('Cookie', cookie);

    expect(response.statusCode).toBe(200);
  });
});

describe('Logout endpoint', () => {
  it('logout', async () => {
    const cookie = await GetCookie();

    expect(cookie).toBeDefined();

    const response = await request(baseUrl)
      .post('/api/v1/auth/logout')
      .set('Cookie', cookie);

    expect(response.statusCode).toBe(200);
  });
});

describe('Delete endpoint', () => {
  it('Delete un-logged user', async () => {
    const response = await request(baseUrl)
      .delete('/api/v1/auth/user');

    expect(response.statusCode).toBe(401);
  });

  it('Delete logged user', async () => {
    const cookie = await GetCookie();
    expect(cookie).toBeDefined();

    const response = await request(baseUrl)
      .delete('/api/v1/auth/user')
      .set('Cookie', cookie);

    expect(response.statusCode).toBe(200);
  });
});
