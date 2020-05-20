export default {
  header: {
    navigation: {
      whitepapers: 'Документация',
      downloads: 'Загрузки',
      explorer: 'Обозреватель',
      blog: 'Блог',
      roadmap: 'План',
    },
  },
  footer: {
    getStarted: 'Начать',
    explore: 'Дополнительно',
    community: 'Сообщество',
    about: 'О проекте',
    wallet_android: 'Android',
    wallet_mac: 'MacOS',
    wallet_win: 'Windows',
    wallet_lin: 'Linux',
    infographics: 'Инфографика',
    whitepapers: 'Документация',
    blockchain: 'Блокчейн Обозреватель',
    blog: 'Блог',
    twitter: 'Twitter',
    reddit: 'Reddit',
    github: 'Github',
    telegram: 'Telegram',
    slack: 'Slack',
    roadmap: 'План',
    skyMessenger: 'Sky-Messenger',
    cxPlayground: 'CX Playground',
    team: 'Команда',
    subscribe: 'Рассылка',
    market: 'Markets',
    bitcoinTalks: 'Bitcointalks ANN',
    instagram: 'Instagram',
    facebook: 'Facebook',
    discord: 'Discord',
    platform: "Платформа",
    explorer: "Explorer",
    price: "CoinPaprika",

  },
  distribution: {
    rate: 'Цена: 1 {coinType} = {rate} MDL',
    inventory: 'Доступно: <strong>{coins} MDL</strong>',
    title: 'MDL Банкомат',
    heading: 'MDL Банкомат',
    headingEnded: 'The previous distribution event finished on',
    ended: `<p>Join the <a href="https://t.me/mdl">MDL Telegram</a>
      or follow the
      <a href="https://twitter.com/mdlproject">MDL Twitter</a>
      to learn when the next event begins.`,
    instructions: `<p>Рыночнкую цену <a href="https://coinmarketcap.com/currencies/mdl/">MDL смотреть на CoinPaprika</a>.</p>

<p>Для преобретения MDL токенов:</p>

<ul>
  <li>Введите ваш MDL адрес</li>
  <li>Вы получите уникальный {coinType} адрес для покупки MDL</li>
  <li>Пошлите {coinType} на полученый адрес</li>
</ul>

<p>Вы можете проверить статус заказа, введя адрес MDL и нажав на <strong>Проверить статус</strong>.</p>
<p>Каждый раз при нажатии на <strong>Получить адрес</strong>, генерируется новый SKY адрес.</p>
<p>Один адрес MDL может иметь не более {max_bound_addrs} прикреплённых адресов Skycoin.</p>
    `,
    statusFor: 'Статус по {mdlAddress}',
    enterAddress: 'Введите адрес MDL',
    getAddress: 'Получить адрес',
    checkStatus: 'Проверить статус',
    loading: 'Загрузка...',
    recAddress: '{coinType} адрес',
    errors: {
      noSkyAddress: 'Пожалуйста введите ваш MDL адрес.',
      coinsSoldOut: 'MDL OTC is currently sold out, check back later.',
    },
    statuses: {
      waiting_deposit: '[tx-{id} {updated}] Ожидаем депозит.',
      waiting_send: '[tx-{id} {updated}] Депозит подтверждён. MDL транзакция поставлена в очередь.',
      waiting_confirm: '[tx-{id} {updated}] MDL транзакция отправлена. Ожидаем подтверждение.',
      done: '[tx-{id} {updated}] Завершена. Проверьте ваш MDL кошелёк.',
    },
  },
};
