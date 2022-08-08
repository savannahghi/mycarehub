BEGIN;
ALTER TABLE questionnaires_questioninputchoice ADD CONSTRAINT  questionnaires_questioninputchoice_choice_key UNIQUE(choice);
ALTER TABLE questionnaires_question ADD CONSTRAINT questionnaires_question_text_key UNIQUE(text);
COMMIT;
