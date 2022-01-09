import Service, { inject as service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import ENV from 'corgi/config/environment';

export default class SessionService extends Service {
  @service store;
  @service session;

  @tracked logged = null;
  @tracked token = null;
  @tracked user = null;
  @tracked userName = null;

  static STORAGE_KEY = 'realworld.ember.token';

  initSession() {
    //let storedToken = this.getStoredToken();
    let accessToken = this.getCookie('access_token');
    if (accessToken) {
      this.token = this.getCookie('access_token');
      return this.fetchUser();
    }
  }

  get isLoggedIn() {
    let accessToken = this.getCookie('access_token');
    if (accessToken) {
      return this.fetchUser();
    }
    return !!accessToken;
  }

  async fetch(url, method = 'GET') {
    let response = await fetch(`${ENV.APP.apiHost}${url}`, {
      method,
      credentials: 'include',
    });
    return await response.json();
  }

  @action
  async register(username, email, password) {
    let user = this.store.createRecord('user', {
      username,
      email,
      password,
    });
    try {
      await user.save();
      this.setToken(user.token);
    } catch (e) {
      // eslint-disable-next-line no-console
      console.error(e);
    } finally {
      this.user = user;
    }
    return user;
  }

  @action
  async logIn(email, password) {
    // @patocallaghan - It would be nice to encapsulate some of this logic in the User model as a `static` class, but unsure how to access container and store from there
    let login = await fetch(`${ENV.APP.apiHost}/auth/password/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        email: email,
        password: password,
      }),
    });

    let response = await login;
    let data = response.json();
    if (response.status !== 200) {
      return { status: response.status, message: data.message };
    }
  }

  @action
  async logOut() {
    this.removeToken();
  }

  async fetchUser() {
    let data = await this.session.fetch('/api/v1/user/me');
    this.store.set('user', data);
    this.user = this.store.get('user');
    return this.user;
  }

  getStoredToken() {
    return localStorage.getItem(SessionService.STORAGE_KEY);
  }

  setToken(token) {
    this.token = token;
    localStorage.setItem(SessionService.STORAGE_KEY, token);
  }

  removeToken() {
    document.cookie = 'access_token=; expires=-1';
    document.cookie = 'refresh_token_id=; expires=-1';
  }

  processLoginErrors(errors) {
    let loginErrors = [];
    let errorKeys = Object.keys(errors);
    errorKeys.forEach((attribute) => {
      errors[attribute].forEach((message) => {
        loginErrors.push(`${attribute} ${message}`);
      });
    });
    return loginErrors;
  }

  getCookie(cookieName) {
    let name = cookieName + '=';
    let decodedCookie = decodeURIComponent(document.cookie);
    let ca = decodedCookie.split(';');
    for (let i = 0; i < ca.length; i++) {
      let c = ca[i];
      while (c.charAt(0) === ' ') {
        c = c.substring(1);
      }
      if (c.indexOf(name) === 0) {
        return c.substring(name.length, c.length);
      }
    }
    return '';
  }
}
