"""
Provide load test for gess node.
"""
from runner.environment import RunnerEnvironment
from runner import runner


if __name__ == '__main__':

    NODES_COUNT = 2
    LOAD_FACTOR = 10

    runner_environment = RunnerEnvironment(nodes_count=NODES_COUNT, load_factor=LOAD_FACTOR)

    try:
        if runner_environment.setup_accounts():
            # Stock run
            runner.run(nodes_count=NODES_COUNT, load_factor=LOAD_FACTOR)
            # Test single run
            # runner.test_run(nodes_count=NODES_COUNT, load_factor=LOAD_FACTOR)

            print('We are done!')
            runner_environment.cleanup()

    except (KeyboardInterrupt, Exception) as exception:
        print('\nCaught an exception. Saving accounts...\n')
        runner_environment.save_accounts()

        raise exception
