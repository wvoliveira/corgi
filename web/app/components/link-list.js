import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency-decorators';
import { inject as service } from '@ember/service';

export default class ArticleListComponent extends Component {
  @service session;
  @service store;
  @tracked links = [];

  constructor() {
    super(...arguments);
    this.loadLinks.perform();
  }

  @task({ restartable: true })
  *loadLinks() {
    let NUMBER_OF_ARTICLES = 10;
    let offset = (parseInt(this.args.page, 10) - 1) * NUMBER_OF_ARTICLES;

    if (this.args.feed === 'your') {
      this.articles = yield this.session.user.fetchFeed(this.args.page);
    } else {
      this.articles = yield this.store.query('article', {
        limit: NUMBER_OF_ARTICLES,
        offset,
        tag: this.args.tag,
      });
    }
  }
}
