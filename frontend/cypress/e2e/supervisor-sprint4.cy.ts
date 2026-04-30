describe('Supervisor Dashboard - Sprint 4 Features', () => {
  beforeEach(() => {
    cy.visit('/home');
    cy.contains('Supervisor').click();
    cy.get('input[name="email"]').type('supervisor2@pandaexpress.com');
    cy.get('input[name="password"]').type('pandasupervisor2');
    cy.contains('Sign In').click();
    cy.url().should('include', '/supervisor-dashboard');
  });

  it('should open the Worker Directory and display company workers', () => {
    cy.contains('.summary-card', 'Worker Directory').click();

    cy.get('.shift-modal.open').should('be.visible');
    cy.get('.shift-modal').contains('Human Resources', { matchCase: false }).should('exist');
    cy.get('th').contains('Name').should('exist');
    cy.get('th').contains('Email').should('exist');

    cy.get('button').contains('Close').click();
    cy.get('.shift-modal.open').should('not.exist');
  });

  it('should unassign a worker from an active shift', () => {
    cy.get('.schedule-grid')
      .contains('Assigned')
      .parent()
      .parent() 
      .find('button')
      .contains('Remove')
      .click();

    cy.get('.mat-mdc-snack-bar-container').contains('Worker unassigned successfully!');
  });

  it('should completely delete an unassigned shift', () => {

    cy.contains('Open Shifts').click();

    cy.get('.shift-modal.open').should('be.visible');

    cy.get('.shift-modal.open')
      .find('.list.compact li')
      .first()
      .find('button')
      .contains('Delete') 
      .click();

    cy.get('.mat-mdc-snack-bar-container').contains('Shift deleted successfully!');
  });
});