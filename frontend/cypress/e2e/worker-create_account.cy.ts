describe('Worker Create Account Flow', () => {
  it('creates worker account using Panda Express', () => {
    cy.intercept('GET', '**/api/companies/').as('getCompanies');
    cy.visit('http://localhost:4200/home');
    cy.contains('Worker').click();
    cy.url().should('include', '/login?role=worker');
    cy.contains('Create Account').click();
    cy.url().should('include', '/register?role=worker');
    cy.wait('@getCompanies');

    // open custom company dropdown
    cy.contains(/Select Company/i).click();

    // choose Panda Express in the overlay list
    cy.get('body')
      .contains('Panda Express', { timeout: 10000 })
      .click();

    cy.get('input[name="name"], input[placeholder*="Name"]').type('raj');
    cy.get('input[name="email"], input[placeholder*="Email"]').type('raj@panda.com');
    cy.get('input[name="password"], input[placeholder*="Password"]').type('123456789');
    cy.get(
      'input[name="number"], input[name="phone"], input[name="phoneNumber"], input[placeholder*="Phone"]',
    ).type('3224502503');
    cy.get(
      'input[name="address"], textarea[name="address"], input[placeholder*="Address"], textarea[placeholder*="Address"]',
    ).type('2400 sw road 13 street Florida');

    cy.contains(/Create Account/i).click();
    cy.contains(/Account created successfully/i, { timeout: 10000 }).should('be.visible');
  });
});