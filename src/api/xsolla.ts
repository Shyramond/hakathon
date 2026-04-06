/**
 * Mock implementation of Xsolla API.
 * In production, you would use actual API endpoints:
 * - Login API: https://login.xsolla.com/api/
 * - Catalog API: https://store.xsolla.com/api/v2/project/{project_id}/
 * - Paystation API: Integrated via SDK or opening URL with token.
 */

export const XsollaAPI = {
  Login: {
    authenticate: async (username: string, _pass: string) => {
      // Mocking Login API
      return new Promise((resolve) => {
        setTimeout(() => {
          resolve({
            token: 'mock_jwt_token',
            user: {
              id: 'user_x',
              name: username,
              email: `${username}@test.com`
            }
          })
        }, 500);
      });
    }
  },
  
  Catalog: {
    getItems: async () => {
      // Mocking Catalog API for virtual items
      return new Promise((resolve) => {
        setTimeout(() => {
          resolve([
            {
              sku: 'stone_pack_grey',
              name: 'Набор серых камней (10 шт)',
              description: 'Базовые материалы для старта ковки.',
              price: 1.99,
              currency: 'USD',
              image: 'grey_stone_bundle'
            },
            {
              sku: 'stone_pack_green',
              name: 'Сундук ученика (Зеленые камни)',
              description: 'Немного удачи не повредит.',
              price: 4.99,
              currency: 'USD',
              image: 'green_stone_bundle'
            },
            {
              sku: 'battle_pass',
              name: 'Сезонный пропуск: Эпоха Огня',
              description: 'Разблокируйте премиальную ленту наград.',
              price: 9.99,
              currency: 'USD',
              image: 'battle_pass'
            }
          ])
        }, 500);
      });
    }
  },

  Paystation: {
    purchaseItem: async (sku: string) => {
      // Mocking Paystation flow (normally redirects to Paystation widget)
      return new Promise((resolve, reject) => {
        // Simulating the iframe/popup interaction
        setTimeout(() => {
          const success = Math.random() > 0.1; // 90% success rate
          if (success) {
            resolve({ status: 'done', sku });
          } else {
            reject(new Error('Payment failed or cancelled'));
          }
        }, 1000);
      });
    }
  }
};
