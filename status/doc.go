/*
Package status is used by Cairo to represent all kinds of errors.
A status value of `status.Success“ represents no error and has an
integer value of zero. All other status values represent an error.

Cairo's error handling is designed to be easy to use and safe.
All major cairo objects retain an error status internally which
can be queried anytime by the users using “*.Status()“ calls.
In the mean time, it is safe to call all cairo functions normally
even if the underlying object is in an error status.
This means that no error handling code is required before or after
each individual cairo function call.
*/
package status
