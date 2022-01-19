// validation documentation: https://github.com/horprogs/Just-validate

// login form validation
const loginValidation = new window.JustValidate('#login-form', {
   errorFieldCssClass: 'is-invalid',
   successFieldCssClass: 'is-valid',
   lockForm: true,
});
loginValidation
   .addField('#username', [
      {rule: 'required', errorMessage: 'Username is required!'},
      {rule: 'email', errorMessage: 'Username must be a valid email address!'}
   ])
   .addField("#password", [
      {rule:'required', errorMessage:'Password is required!'},
      {rule: 'password'}
   ])
   .onSuccess((e) => {
      console.log('validation passed!', e);
      // remember this!
      e.submitter.form.submit();
   });
