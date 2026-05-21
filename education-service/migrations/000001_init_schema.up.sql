CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE subjects (
                          id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                          title VARCHAR(255) NOT NULL,
                          credits INT NOT NULL,
                          teacher_id UUID NOT NULL
);

CREATE TABLE teachers (
                          id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                          name VARCHAR(255) NOT NULL,
                          department VARCHAR(255) NOT NULL,
                          email VARCHAR(255),
                          office VARCHAR(100)
);

CREATE TABLE schedule (
                          id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                          subject_id UUID REFERENCES subjects(id) ON DELETE CASCADE,
                          room VARCHAR(50) NOT NULL,
                          start_time TIMESTAMP NOT NULL,
                          day_of_week INT NOT NULL
);

CREATE TABLE grades (
                        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                        student_id UUID NOT NULL,
                        subject_id UUID REFERENCES subjects(id) ON DELETE CASCADE,
                        score FLOAT NOT NULL,
                        type VARCHAR(50) NOT NULL
);

CREATE TABLE attendance (
                            id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                            student_id UUID NOT NULL,
                            schedule_id UUID REFERENCES schedule(id) ON DELETE CASCADE,
                            status VARCHAR(20) NOT NULL,
                            date DATE NOT NULL
);

CREATE TABLE student_courses (
                                 student_id UUID NOT NULL,
                                 subject_id UUID REFERENCES subjects(id) ON DELETE CASCADE,
                                 PRIMARY KEY (student_id, subject_id)
);

-- ===== SEED DATA =====

INSERT INTO teachers (id, name, department, email, office) VALUES
                                                               ('aaaaaaaa-0001-0001-0001-000000000001', 'Dr. Aizat Bekova',    'Computer Science',  'bekova@aitu.edu.kz',    'B-201'),
                                                               ('aaaaaaaa-0001-0001-0001-000000000002', 'Prof. Damir Seitkali','Mathematics',        'seitkali@aitu.edu.kz',  'A-305'),
                                                               ('aaaaaaaa-0001-0001-0001-000000000003', 'Dr. Alina Nurova',    'Software Engineering','nurova@aitu.edu.kz',   'C-110'),
                                                               ('aaaaaaaa-0001-0001-0001-000000000004', 'Prof. Ruslan Akhmet', 'Data Science',       'akhmet@aitu.edu.kz',    'B-402'),
                                                               ('aaaaaaaa-0001-0001-0001-000000000005', 'Dr. Saltanat Yesova', 'Cybersecurity',      'yesova@aitu.edu.kz',    'A-108');

INSERT INTO subjects (id, title, credits, teacher_id) VALUES
                                                          ('bbbbbbbb-0001-0001-0001-000000000001', 'Algorithms & Data Structures', 5, 'aaaaaaaa-0001-0001-0001-000000000001'),
                                                          ('bbbbbbbb-0001-0001-0001-000000000002', 'Calculus II',                  4, 'aaaaaaaa-0001-0001-0001-000000000002'),
                                                          ('bbbbbbbb-0001-0001-0001-000000000003', 'Software Architecture',        4, 'aaaaaaaa-0001-0001-0001-000000000003'),
                                                          ('bbbbbbbb-0001-0001-0001-000000000004', 'Machine Learning',             5, 'aaaaaaaa-0001-0001-0001-000000000004'),
                                                          ('bbbbbbbb-0001-0001-0001-000000000005', 'Network Security',             3, 'aaaaaaaa-0001-0001-0001-000000000005');

INSERT INTO schedule (id, subject_id, room, start_time, day_of_week) VALUES
                                                                         ('cccccccc-0001-0001-0001-000000000001', 'bbbbbbbb-0001-0001-0001-000000000001', '301A', '2025-01-01 09:00:00', 1),
                                                                         ('cccccccc-0001-0001-0001-000000000002', 'bbbbbbbb-0001-0001-0001-000000000002', '205B', '2025-01-01 11:00:00', 1),
                                                                         ('cccccccc-0001-0001-0001-000000000003', 'bbbbbbbb-0001-0001-0001-000000000003', '401C', '2025-01-01 14:00:00', 2),
                                                                         ('cccccccc-0001-0001-0001-000000000004', 'bbbbbbbb-0001-0001-0001-000000000004', '102A', '2025-01-01 09:00:00', 3),
                                                                         ('cccccccc-0001-0001-0001-000000000005', 'bbbbbbbb-0001-0001-0001-000000000005', '303B', '2025-01-01 16:00:00', 4);