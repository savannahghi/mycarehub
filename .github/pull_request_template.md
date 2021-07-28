# Review Checklist

[![Bugs](https://img.shields.io/badge/BugFixes-0-red.svg?style=for-the-badge)](https://shields.io/)
[![Features](https://img.shields.io/badge/Features-0-orange?style=for-the-badge)](https://shields.io/)
[![Tests](https://img.shields.io/badge/Tests-0-success.svg?style=for-the-badge)](https://shields.io/)

#### Summary*

- [ ] fix: Bug fix for `Add the items tackled here as checklists`
- [x] feat: Completed task
- [ ] chore: Incomplete task
  - [ ] test: Sub-task 1


#### Structure*

- [ ] The Pull Request has a `proper` title that conforms to our [MR title standards](https://gist.github.com/mikepea/863f63d6e37281e329f8)
- [ ] The Pull Request has one commit, and if there are more than one, they should be squashed
- [ ] The commit should have a proper title and a short description
- [ ] The commit must be signed off
- [ ] Unused imports are not present
- [ ] Dead code is not present
- [ ] Ensure dry running library to confirm changes

#### Tests

- [ ] Proper and high quality unit, integration and acceptance(if applicable) tests have been written
- [ ] The coverage threshold should not be lowered

### Sign off*

- [ ] All comments have been resolved by the reviewers
- [ ] Approved by Czar {replace_with_czar_name}
- [ ] Signed off by second reviewer {replace_with_name}
- [ ] Ensure all checks are done before merging :warning:
- [ ] All PRs needs to be signed before merging :warning:

#### N/B:

- Add a checklist if more than one item is done.
- Add screenshots and/or images where necessary
- Indicate any breakages caused in the UI :exclamation:
- Where necessary, indicate which issue the Pull Request solves (Closes #)
- Any new files are updated in the folder structure in the [README](../../README.md)