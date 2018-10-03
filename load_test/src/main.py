"""
Provide load test for gess node.
"""
from runner.environment import RunnerEnvironment
from runner import runner


if __name__ == '__main__':

    NODES_COUNT = 2
    LOAD_FACTOR = 1000

    runner_environment = RunnerEnvironment(nodes_count=NODES_COUNT, load_factor=LOAD_FACTOR)

    try:
        if runner_environment.setup_accounts():
            runner.run(nodes_count=NODES_COUNT, load_factor=LOAD_FACTOR)
            # runner.test_run(nodes_count=NODES_COUNT, load_factor=LOAD_FACTOR)
            print('We are done!')
            runner_environment.cleanup()

    except Exception as exception:
        runner_environment.save_accounts()

        raise exception
